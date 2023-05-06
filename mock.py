
import random
import psycopg2
from datetime import date, datetime, time, timedelta
import nltk
import ssl

# try:
#     _create_unverified_https_context = ssl._create_unverified_context
# except AttributeError:
#     pass
# else:
#     ssl._create_default_https_context = _create_unverified_https_context

# Set up database connection
conn = psycopg2.connect(
    host="localhost",
    database="aimet",
    user="aimet",
    password="aimetpassword"
)

# nltk.download('brown')
# Define parameters for generating events

min_start_date = date(2010, 1, 1)      # Start date for generating events
max_end_date = date(2040, 12, 31)      # End date for generating events
min_start_time = time(1, 0)            # Start time for events
max_end_time = time(20, 0)             # End time for events
max_event_duration = 60                # Duration of each event in minutes

# Define function to generate random datetime within a date range
def random_datetime(start_date, end_date, start_time, end_time):

    # print(start_time.hour, end_time.hour)
    random_date = random.choice([d for d in range((end_date - start_date).days + 1)])
    random_time = time(
        random.randint(start_time.hour, end_time.hour),
        random.randint(0, 55)  # Round minutes to nearest 15-minute interval
    )
    return datetime.combine(start_date + timedelta(days=random_date), random_time)



def generate_random_phrase():
    # Download the required nltk data
    

    # Load the list of words
    words = nltk.corpus.brown.words()
    
    title_list = []
    for i in range(random.randint(2, 4)):
        title_list.append(random.choice(words))
    # Choose a random adjective and noun
    
    # Generate a random number between 1 and 999
    number = random.randint(1, 999)

    # Return the phrase with the number appended
    title = " ".join(title_list)
    return f"{title} {number}"

# Generate and insert events into the database

cur = conn.cursor()
for j in range(20000):
    for k in range(300):
        # Generate random event properties
        title = generate_random_phrase()
        event_datetime = random_datetime(min_start_date, max_end_date, min_start_time, max_end_time)
        start_time = event_datetime.time()
        end_time = (event_datetime + timedelta(minutes=random.randint(40, 300))).time()
        
        try:
            # Insert event into the database
            cur.execute("""
                INSERT INTO events (title, event_date, start_time, end_time)
                VALUES (%s, %s, %s, %s)
            """, (title, event_datetime.date(), start_time, end_time))
            conn.commit()
            
            break  # Break out of while loop if event is inserted successfully
        except psycopg2.errors.UniqueViolation:
            # If event overlaps with existing event, generate a new event
            continue
        except psycopg2.errors.InFailedSqlTransaction:
            # If transaction is aborted, rollback and start a new transaction
            conn.rollback()
            continue
        except psycopg2.errors.RaiseException as e:
        # If event violates overlapping constraint, print error message and skip event
            conn.rollback()
            continue
        except psycopg2.errors.CheckViolation as e:
        # If event violates overlapping constraint, print error message and skip event
            conn.rollback()
            continue
# Close database connection
conn.close()